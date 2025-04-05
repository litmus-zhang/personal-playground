import pytest
import pandas as pd
import numpy as np
from main import (  # Replace 'your_module' with your actual module name
    get_forex_data,
    test_cointegration,
    calculate_hedge_ratio,
    generate_signals,
    calculate_position_size,
    backtest_strategy
)

# Fixture for test data
@pytest.fixture
def test_data():
    """Fixture that provides test data for cointegrated pairs"""
    # Use AUD/USD and NZD/USD as they're typically cointegrated
    data = get_forex_data('AUDUSD=X', 'NZDUSD=X')
    return data

def test_get_forex_data(test_data):
    """Test data fetching function"""
    assert isinstance(test_data, pd.DataFrame)
    assert len(test_data) > 100  # Should have sufficient data points
    assert set(test_data.columns) == {'audusd', 'nzdusd'}
    assert not test_data.isnull().any().any()  # No NaN values

def test_cointegration_test(test_data):
    """Test cointegration testing function"""
    p_value, is_cointegrated = test_cointegration(test_data['audusd'], test_data['nzdusd'])
    
    assert isinstance(p_value, float)
    assert isinstance(is_cointegrated, bool)
    assert 0 <= p_value <= 1
    # For these typically cointegrated pairs, we expect True
    assert is_cointegrated == True

def test_hedge_ratio_calculation(test_data):
    """Test hedge ratio calculation"""
    model, hedge_ratio = calculate_hedge_ratio(test_data['nzdusd'], test_data['audusd'])
    
    assert hasattr(model, 'params')  # Verify it's a proper statsmodels result
    assert isinstance(hedge_ratio, float)
    # Hedge ratio should be positive and reasonable for these pairs
    assert 0.5 < hedge_ratio < 2.0

def test_signal_generation(test_data):
    """Test signal generation function"""
    # First calculate spread
    _, hedge_ratio = calculate_hedge_ratio(test_data['nzdusd'], test_data['audusd'])
    spread = test_data['nzdusd'] - hedge_ratio * test_data['audusd']
    
    signals = generate_signals(spread)
    
    assert isinstance(signals, pd.DataFrame)
    assert {'zscore', 'position'}.issubset(signals.columns)
    assert signals['position'].isin([-1, 0, 1]).all()
    # Verify we have some signals (not all zeros)
    assert signals['position'].abs().sum() > 0

def test_position_sizing(test_data):
    """Test position sizing with risk management"""
    # First generate signals
    _, hedge_ratio = calculate_hedge_ratio(test_data['nzdusd'], test_data['audusd'])
    spread = test_data['nzdusd'] - hedge_ratio * test_data['audusd']
    signals = generate_signals(spread)
    
    # Test position sizing
    sized_positions = calculate_position_size(signals, spread)
    
    assert isinstance(sized_positions, pd.DataFrame)
    assert {'position_size', 'stop_loss'}.issubset(sized_positions.columns)
    # Position sizes should be positive when in a trade
    assert (sized_positions.loc[sized_positions['position'] != 0, 'position_size'] > 0).all()
    # Stop loss should be set when in a trade
    assert sized_positions.loc[sized_positions['position'] != 0, 'stop_loss'].notna().all()

def test_backtest(test_data):
    """Test complete backtest workflow"""
    # Run full workflow
    _, hedge_ratio = calculate_hedge_ratio(test_data['nzdusd'], test_data['audusd'])
    spread = test_data['nzdusd'] - hedge_ratio * test_data['audusd']
    signals = generate_signals(spread)
    sized_positions = calculate_position_size(signals, spread)
    results, metrics = backtest_strategy(sized_positions, test_data, hedge_ratio)
    
    # Test results
    assert isinstance(results, pd.DataFrame)
    assert isinstance(metrics, dict)
    assert {'strategy_returns', 'cumulative_returns'}.issubset(results.columns)
    
    # Test metrics
    required_metrics = {
        'annualized_return',
        'annualized_volatility',
        'sharpe_ratio',
        'max_drawdown',
        'win_rate',
        'total_trades'
    }
    assert required_metrics.issubset(metrics.keys())
    
    # Basic sanity checks on metrics
    assert metrics['total_trades'] > 0
    assert -1 <= metrics['max_drawdown'] <= 0
    assert 0 <= metrics['win_rate'] <= 1

def test_end_to_end_workflow():
    """Complete end-to-end test of the entire system"""
    # 1. Get data
    data = get_forex_data('AUDUSD=X', 'NZDUSD=X', period='1y')
    
    # 2. Test cointegration
    p_value, is_cointegrated = test_cointegration(data['audusd'], data['nzdusd'])
    assert is_cointegrated, "Test pairs should be cointegrated"
    
    # 3. Calculate hedge ratio
    model, hedge_ratio = calculate_hedge_ratio(data['nzdusd'], data['audusd'])
    assert isinstance(hedge_ratio, float)
    
    # 4. Calculate spread
    spread = data['nzdusd'] - hedge_ratio * data['audusd']
    
    # 5. Generate signals
    signals = generate_signals(spread)
    assert not signals.empty
    
    # 6. Calculate position sizes
    sized_positions = calculate_position_size(signals, spread)
    assert not sized_positions.empty
    
    # 7. Backtest strategy
    results, metrics = backtest_strategy(sized_positions, data, hedge_ratio)
    assert not results.empty
    assert isinstance(metrics, dict)
    
    # Verify we have some returns data
    assert results['strategy_returns'].abs().sum() > 0