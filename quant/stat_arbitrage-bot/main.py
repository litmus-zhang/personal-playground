import yfinance as yf
import pandas as pd
import matplotlib.pyplot as plt
from statsmodels.tsa.stattools import coint
import statsmodels.api as sm
import numpy as np

import warnings


warnings.filterwarnings('ignore', category=pd.errors.SettingWithCopyWarning)

# Fetch data with proper tickers and ensure alignment
def get_forex_data(pair1, pair2, period="5y", interval="1d"):
    """
    Fetch and align forex data from yfinance
    Returns merged DataFrame with closing prices
    """
    
    df1 = yf.download(pair1, period=period, interval=interval)
    df2 = yf.download(pair2, period=period, interval=interval)
    
    # Clean and rename columns
    df1 = df1[['Close']].rename(columns={'Close': pair1.replace('=X', '').lower()})
    df2 = df2[['Close']].rename(columns={'Close': pair2.replace('=X', '').lower()})
    
    # Merge and forward fill
    merged = pd.merge(df1, df2, left_index=True, right_index=True, how='outer')
    merged = merged.ffill().dropna()
    
    return merged


def test_cointegration(series1, series2):
    """
    Test for cointegration between two price series
    Returns p-value and whether they're cointegrated at 95% confidence
    """
    score, p_value, _ = coint(series1, series2)
    is_cointegrated = p_value < 0.05
    return p_value, is_cointegrated


def calculate_hedge_ratio(y, x):
    """
    Calculate hedge ratio using OLS regression
    Returns model and hedge ratio
    """
    x = sm.add_constant(x)  # Add constant for intercept
    model = sm.OLS(y, x).fit()
    hedge_ratio = model.params.iloc[1]  # Use iloc to avoid deprecation warning
    return model, hedge_ratio


def generate_signals(spread, entry_z=1.5, exit_z=0.5, lookback=20):
    """
    Generate trading signals based on z-score of spread
    Parameters:
        spread (pd.Series): The calculated spread series
        entry_z (float): Z-score threshold for entry
        exit_z (float): Z-score threshold for exit
        lookback (int): Window for rolling z-score calculation
    Returns:
        pd.DataFrame with signals and positions
    """
    # Calculate rolling z-score
    spread_mean = spread.rolling(lookback).mean()
    spread_std = spread.rolling(lookback).std()
    zscore = (spread - spread_mean) / spread_std
    
    # Initialize signals
    signals = pd.DataFrame(index=spread.index)
    signals['zscore'] = zscore
    signals['position'] = 0  # 0: flat, 1: long spread, -1: short spread
    
    # Generate signals
    position = 0
    for i in range(len(signals)):
        # Entry conditions
        if position == 0:
            if signals['zscore'].iloc[i] > entry_z:
                signals['position'].iloc[i] = -1  # Short spread
                position = -1
            elif signals['zscore'].iloc[i] < -entry_z:
                signals['position'].iloc[i] = 1  # Long spread
                position = 1
        # Exit conditions
        elif position != 0:
            if (position == 1 and signals['zscore'].iloc[i] > -exit_z) or \
               (position == -1 and signals['zscore'].iloc[i] < exit_z):
                signals['position'].iloc[i] = 0  # Close position
                position = 0
            else:
                signals['position'].iloc[i] = position  # Maintain position
    
    return signals


def calculate_position_size(signals, spread, capital=10000, risk_pct=0.01, stop_loss_z=2.0):
    """
    Calculate position sizes with risk management - FIXED VERSION
    """
    # Create a copy to avoid modifying the original
    df = signals.copy()
    df['spread'] = spread
    
    # Calculate rolling stats for stop loss
    lookback = 20
    spread_mean = spread.rolling(lookback).mean()
    spread_std = spread.rolling(lookback).std()
    
    # Initialize new columns
    df['position_size'] = 0.0
    df['stop_loss'] = np.nan
    
    # Get index as list for sequential access
    index_list = df.index.tolist()
    
    for i, current_index in enumerate(index_list):
        # Get previous index if not first element
        prev_index = index_list[i-1] if i > 0 else None
        
        # Entry conditions
        if df.at[current_index, 'position'] != 0:
            if i == 0 or df.at[current_index, 'position'] != df.at[prev_index, 'position']:
                # Calculate stop loss
                if df.at[current_index, 'position'] == 1:  # Long spread
                    stop_loss = spread_mean.at[current_index] - stop_loss_z * spread_std.at[current_index]
                else:  # Short spread
                    stop_loss = spread_mean.at[current_index] + stop_loss_z * spread_std.at[current_index]
                
                # Calculate position size
                risk_amount = capital * risk_pct
                price_diff = abs(df.at[current_index, 'spread'] - stop_loss)
                position_size = risk_amount / price_diff if price_diff > 0 else 0
                
                # Assign values
                df.at[current_index, 'stop_loss'] = stop_loss
                df.at[current_index, 'position_size'] = position_size
                
        # Exit conditions
        elif i > 0 and df.at[prev_index, 'position'] != 0:
            # Clear values when position closed
            df.at[current_index, 'stop_loss'] = np.nan
            df.at[current_index, 'position_size'] = 0.0
            
        # Maintain position conditions
        elif i > 0:
            # Carry forward existing values
            df.at[current_index, 'stop_loss'] = df.at[prev_index, 'stop_loss']
            df.at[current_index, 'position_size'] = df.at[prev_index, 'position_size']
    
    return df


def backtest_strategy(signals_with_sizes, prices, hedge_ratio):
    """
    Backtest the strategy - FIXED VERSION
    """
    # Make sure we're working with proper DataFrames
    df = signals_with_sizes.copy()
    prices_df = prices.copy()
    
    # Ensure we're only using the close prices if prices has multiple columns
    if isinstance(prices_df, pd.DataFrame) and len(prices_df.columns) > 1:
        prices_df = prices_df[['close_eur', 'close_gbp']]  # Adjust column names as needed
    
    # Join using merge instead of join for better control
    df = pd.merge(df, prices_df, left_index=True, right_index=True, how='left')
    
    # Calculate daily returns
    price_cols = [col for col in df.columns if col.startswith('close_')]
    for i, col in enumerate(price_cols, 1):
        df[f'price{i}_returns'] = df[col].pct_change()
    
    # Calculate strategy returns
    df['strategy_returns'] = 0.0
    for i in range(1, len(df)):
        if df['position'].iloc[i-1] == 1:  # Long spread
            df.loc[df.index[i], 'strategy_returns'] = (
                df['position_size'].iloc[i-1] * 
                (df['price2_returns'].iloc[i] - hedge_ratio * df['price1_returns'].iloc[i])
            )
        elif df['position'].iloc[i-1] == -1:  # Short spread
            df.loc[df.index[i], 'strategy_returns'] = (
                df['position_size'].iloc[i-1] * 
                (-df['price2_returns'].iloc[i] + hedge_ratio * df['price1_returns'].iloc[i])
            )
    
    # Calculate cumulative returns
    df['cumulative_returns'] = (1 + df['strategy_returns']).cumprod()
    
    # Calculate performance metrics
    annualized_return = df['strategy_returns'].mean() * 252
    annualized_vol = df['strategy_returns'].std() * np.sqrt(252)
    sharpe_ratio = annualized_return / annualized_vol if annualized_vol != 0 else 0
    
    max_drawdown = (df['cumulative_returns'] / df['cumulative_returns'].cummax() - 1).min()
    win_rate = (df['strategy_returns'] > 0).mean()
    
    metrics = {
        'annualized_return': annualized_return,
        'annualized_volatility': annualized_vol,
        'sharpe_ratio': sharpe_ratio,
        'max_drawdown': max_drawdown,
        'win_rate': win_rate,
        'total_trades': (df['position'].diff() != 0).sum()
    }
    
    return df, metrics


# Test our pairs
forex_pairs = ('AUDUSD=X', 'NZDUSD=X')

# 1. Get data
prices = get_forex_data(forex_pairs[0], forex_pairs[1])

# 2. Test cointegration
p_value, is_cointegrated = test_cointegration(prices.iloc[:, 0], prices.iloc[:, 1])
if not is_cointegrated:
    raise ValueError("Selected pairs are not cointegrated")

# 3. Calculate hedge ratio
model, hedge_ratio = calculate_hedge_ratio(prices.iloc[:, 1], prices.iloc[:, 0])
spread = prices.iloc[:, 1] - hedge_ratio * prices.iloc[:, 0]

# 4. Generate signals
signals = generate_signals(spread, entry_z=1.5, exit_z=0.5)

# 5. Calculate position sizes with risk management
signals_with_sizes = calculate_position_size(signals, spread)

# 6. Backtest strategy
results, metrics = backtest_strategy(signals_with_sizes, prices, hedge_ratio)

# 7. Visualize results
import matplotlib.pyplot as plt

fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 8), sharex=True)

# Plot spread and signals
ax1.plot(results.index, results['spread'], label='Spread')
ax1.plot(results.index, results['zscore'], label='Z-Score', alpha=0.5)
ax1.axhline(0, color='black', linestyle='--')
ax1.set_ylabel('Spread Value')
ax1.legend()

# Plot cumulative returns
ax2.plot(results.index, results['cumulative_returns'], label='Strategy')
ax2.set_ylabel('Cumulative Returns')
ax2.legend()

plt.tight_layout()
plt.savefig("Strategy_returns.svg")
plt.show()

# Print metrics
print("\nStrategy Performance Metrics:")
for k, v in metrics.items():
    print(f"{k.replace('_', ' ').title():<20}: {v:.4f}")