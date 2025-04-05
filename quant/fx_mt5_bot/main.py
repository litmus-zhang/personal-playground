import pandas as pd
import numpy as np
import os
import logging
import time
import json
import pickle
import warnings
import requests
import MetaTrader5 as mt5
import tensorflow as tf
import ta
import matplotlib.pyplot as plt
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sklearn.model_selection import TimeSeriesSplit
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score, mean_absolute_percentage_error
from tensorflow.keras.models import Sequential, load_model, save_model
from tensorflow.keras.layers import LSTM, Dense, Dropout, BatchNormalization
from tensorflow.keras.callbacks import EarlyStopping, ModelCheckpoint, ReduceLROnPlateau
from tensorflow.keras.optimizers import Adam
from statsmodels.tsa.stattools import adfuller
from datetime import datetime, timedelta
import sys


if sys.stdout.encoding != 'utf-8':
    sys.stdout.reconfigure(encoding='utf-8')

warnings.filterwarnings('ignore')

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("forex_trading_bot.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("forex_trading_bot")



class MT5Trader:
    """
    Class to handle trading operations using MetaTrader 5.
    """
    def __init__(self, symbol, lot_size=0.01, tp_multiplier=1.5, sl_multiplier=1.0, 
                 trailing_start_pct=0.3, trailing_step_pct=0.1, risk_per_trade_pct=1.0):
        """
        Initialize the trader.
        
        Parameters:
        -----------
        symbol : str
            Trading pair symbol (e.g., 'EURUSD')
        lot_size : float
            Base lot size for trading
        tp_multiplier : float
            Multiplier for Take Profit
        sl_multiplier : float
            Multiplier for Stop Loss
        trailing_start_pct : float
            Percentage of TP when trailing should activate
        trailing_step_pct : float
            Percentage for trailing step size
        risk_per_trade_pct : float
            Percentage of account balance to risk per trade
        """
    
    def start(self, run_once=False):
        """
        Start the trading bot.
        
        Parameters:
        -----------
        run_once : bool
            If True, run a single trading cycle. If False, run continuously.
        """
        logger.info("Starting trading bot...")
        self.send_telegram_message(f"🚀 Forex trading bot starting for {self.symbol} ({self.timeframe_str})")
        
        if run_once:
            logger.info("Running single trading cycle...")
            return self.run_trading_cycle()
        
        # Run continuously
        logger.info(f"Bot started in continuous mode. Timeframe: {self.timeframe_str}")
        
        # Determine sleep time based on timeframe
        timeframe_minutes = {
            'M1': 1,
            'M5': 5,
            'M15': 15,
            'M30': 30,
            'H1': 60,
            'H4': 240,
            'D1': 1440
        }
        
        sleep_minutes = timeframe_minutes.get(self.timeframe_str, 60)
        
        try:
            while True:
                # Get current time
                now = datetime.now()
                
                # Only trade during weekdays (forex market closed on weekends)
                if now.weekday() >= 5:  # 5=Saturday, 6=Sunday
                    logger.info("Weekend detected. Forex market closed. Waiting...")
                    
                    # Calculate time until Monday
                    days_until_monday = (7 - now.weekday()) % 7
                    if days_until_monday == 0:
                        days_until_monday = 1  # Safety check
                    
                    # Sleep until Monday 00:00 + 1 hour (to ensure markets are open)
                    next_run = now.replace(hour=1, minute=0, second=0) + timedelta(days=days_until_monday)
                    sleep_seconds = (next_run - now).total_seconds()
                    
                    logger.info(f"Sleeping until Monday {next_run} ({sleep_seconds/3600:.1f} hours)")
                    self.send_telegram_message(f"💤 Weekend detected. Bot will resume on Monday.")
                    
                    time.sleep(min(sleep_seconds, 86400))  # Sleep max 24 hours at a time
                    continue
                
                # Run trading cycle
                cycle_result = self.run_trading_cycle()
                
                if cycle_result.get('status') == 'error':
                    logger.error(f"Trading cycle error: {cycle_result.get('message')}")
                
                # Sleep until the next candle
                sleep_seconds = sleep_minutes * 60
                logger.info(f"Sleeping for {sleep_minutes} minutes until next cycle...")
                time.sleep(sleep_seconds)

        except Exception  as e:
            logger.error(f"Error sending Telegram message: {e}")


    def send_telegram_message(self, message):
        """Send notification via Telegram."""
        if not self.telegram_bot_token or not self.chat_id:
            return False
        
        try:
            url = f"https://api.telegram.org/bot{self.telegram_bot_token}/sendMessage"
            payload = {
                "chat_id": self.chat_id,
                "text": message
            }
            response = requests.post(url, json=payload, timeout=10)
            logger.info(f"Telegram message sent: {message}")
            return response.status_code == 200
        except Exception as e:
            logger.error(f"Error sending Telegram message: {e}")
            return False

    
    def build_model(self, input_shape):
        """
        Build a simplified LSTM model focused on predicting price changes.
        """
        try:
            # Create a simpler, more efficient model
            model = Sequential([
                    # Capa de entrada con mayor capacidad
                    LSTM(128, input_shape=input_shape, return_sequences=True, 
                        recurrent_dropout=0.2),
                    BatchNormalization(),
                    Dropout(0.3),
                    
                    # Capas intermedias más robustas
                    LSTM(64, return_sequences=True),
                    BatchNormalization(),
                    Dropout(0.3),
                    
                    # Capa final
                    LSTM(32, return_sequences=False),
                    BatchNormalization(),
                    Dropout(0.2),
                    
                    # Capas densas para mejor aprendizaje no lineal
                    Dense(16, activation='relu'),
                    BatchNormalization(),
                    
                    # Salida: podría ser múltiple para diferentes horizontes temporales
                    Dense(1)
                ])
            
            # Compile with Adam optimizer and learning rate scheduling
            optimizer = Adam(learning_rate=0.001)
            model.compile(optimizer=optimizer, loss='mse', metrics=['mae'])
            
            return model 
        except Exception as e:
            logger.error(f"Error building model: {e}")

    def _plot_predictions(self, y_true, y_pred, timestamp):
        """
        Plot predictions vs actual values.
        
        Parameters:
        -----------
        y_true : numpy.ndarray
            True values
        y_pred : numpy.ndarray
            Predicted values
        timestamp : str
            Timestamp for file naming
        """
        try:
            plt.figure(figsize=(14, 8))
            
            # Plot time series of predictions vs actual
            plt.subplot(2, 1, 1)
            plt.plot(y_true, label='Actual', color='blue', alpha=0.8)
            plt.plot(y_pred, label='Predicted', color='red', linestyle='--', alpha=0.8)
            plt.title(f'{self.symbol} - Actual vs Predicted ({self.timeframe_str})', fontsize=14)
            plt.ylabel('Price Change %', fontsize=12)
            plt.legend(fontsize=12)
            plt.grid(True, alpha=0.3)
            
            # Add metrics to the plot
            mse = mean_squared_error(y_true, y_pred)
            rmse = np.sqrt(mse)
            mae = mean_absolute_error(y_true, y_pred)
            r2 = r2_score(y_true, y_pred)
            
            metrics_text = f'RMSE: {rmse:.6f}\nMAE: {mae:.6f}\nR²: {r2:.6f}'
            plt.annotate(metrics_text, xy=(0.02, 0.85), xycoords='axes fraction', 
                    bbox=dict(boxstyle="round,pad=0.3", fc="white", ec="gray", alpha=0.8))
            
            # Plot correlation scatter plot
            plt.subplot(2, 1, 2)
            plt.scatter(y_true, y_pred, alpha=0.5)
            
            # Add correlation line
            min_val = min(np.min(y_true), np.min(y_pred))
            max_val = max(np.max(y_true), np.max(y_pred))
            plt.plot([min_val, max_val], [min_val, max_val], 'k--')
            
            plt.title('Prediction Correlation', fontsize=14)
            plt.xlabel('Actual', fontsize=12)
            plt.ylabel('Predicted', fontsize=12)
            plt.grid(True, alpha=0.3)
            
            # Add correlation coefficient
            corr = np.corrcoef(y_true.flatten(), y_pred.flatten())[0, 1]
            plt.annotate(f'Correlation: {corr:.4f}', xy=(0.02, 0.9), xycoords='axes fraction',
                    bbox=dict(boxstyle="round,pad=0.3", fc="white", ec="gray", alpha=0.8))
            
            # Save the plot
            plt.tight_layout()
            pred_file = f"{self.plots_dir}/{self.symbol_file}_{self.timeframe_str}_predictions_{timestamp}.png"
            plt.savefig(pred_file, dpi=300)
            plt.close()
        
        except Exception as e:
            logger.error(f"Error plotting prediction: {e}")
        


    def make_trading_decision(self, prediction):
        """
        Make a trading decision based on predicted price change.
        
        Parameters:
        -----------
        prediction : dict
            Prediction data with forecasted values
            
        Returns:
        --------
        dict
            Trading decision
        """
        if prediction is None:
            logger.warning("No prediction available for trading decision")
            return {'type': 'hold', 'confidence': 0, 'reasons': ['No prediction available']}
        
        # Default decision
        decision = {
            'type': 'hold',
            'confidence': 0.5,
            'reasons': []
        }
        
        try:
            # Extract prediction values
            predicted_pct_change = prediction['predicted_pct_change']
            confidence = prediction['confidence']
            market_state = prediction['market_state']
            
            # Calculate decision confidence based on prediction confidence
            # and trend/volatility context
            decision_confidence = confidence
            
            # Adjust based on market state
            if market_state['trend'] == 'up' and predicted_pct_change > 0:
                # Positive prediction in uptrend increases confidence
                decision_confidence = min(0.95, decision_confidence * 1.2)
                decision['reasons'].append(f"Aligned with uptrend")
            elif market_state['trend'] == 'down' and predicted_pct_change < 0:
                # Negative prediction in downtrend increases confidence
                decision_confidence = min(0.95, decision_confidence * 1.2)
                decision['reasons'].append(f"Aligned with downtrend")
        
        except Exception as e:
            logger.error(f"Error making  trading decision: {e}")


    def calculate_position_size(self, entry_price, stop_loss):
        """
        Calculate position size based on risk parameters.
        
        Parameters:
        -----------
        entry_price : float
            Entry price for the position
        stop_loss : float
            Stop loss price level
                
        Returns:
        --------
        float
            Position size in lots
        """
        try:
            # Update account information
            if not self._update_account_info():
                return self.base_lot_size
            
            # Calculate position size based on risk
            risk_amount = self.balance * (self.risk_per_trade_pct / 100)
            
            # Calculate stop loss distance in pips
            stop_distance = abs(entry_price - stop_loss) / self.pip_value
            
            # Calculate lot size based on risk per pip
            if stop_distance > 0:
                # Rule of thumb: 0.1 lot = $1 per pip for most major pairs
                risk_per_pip = risk_amount / stop_distance
                position_size = risk_per_pip / 10  # Approximate conversion to lots
            else:
                position_size = self.base_lot_size
            
            # Adjust to allowed lot size
            position_size = max(self.min_lot, min(self.max_lot, position_size))
            position_size = round(position_size / self.lot_step) * self.lot_step
            
            logger.info(f"Calculated position size: {position_size} lots "
                    f"(Risk: ${risk_amount:.2f}, Stop distance: {stop_distance:.1f} pips)")
            
            return position_size

        except Exception as e:
            logger.error(f"Error calculating position size: {e}")
        
    
    def process_trailing_stops(self):
        """
        Process trailing stops for all active positions.
        
        Returns:
        --------
        int
            Number of positions modified
        """
        try:
            modified_count = 0
            
            # Get all current positions
            positions = mt5.positions_get(symbol=self.symbol)
            if positions is None:
                logger.warning("No positions to process for trailing stop")
                return 0
            
            for position in positions:
                position_ticket = position.ticket
                
                # Skip if not in our tracked positions
                if position_ticket not in self.active_positions:
                    continue
                
                position_data = self.active_positions[position_ticket]
                position_type = position_data['type']
                current_sl = position_data['current_sl']
                original_tp = position_data['original_tp']
                
                # Get current price
                bid, ask = self.get_current_price()
                if bid is None or ask is None:
                    continue
                
                # Calculate current profit
                current_price = bid if position_type == 'buy' else ask
                entry_price = position_data['entry_price']
                
                # Calculate price movement and TP distance
                if position_type == 'buy':
                    price_movement = current_price - entry_price
                    tp_distance = original_tp - entry_price if original_tp else 0
                else:  # sell
                    price_movement = entry_price - current_price
                    tp_distance = entry_price - original_tp if original_tp else 0
                
                # Calculate profit percentage of TP
                profit_pct_of_tp = (price_movement / tp_distance) * 100 if tp_distance > 0 else 0
                
                # Trailing stop logic implementation...
        
        except Exception as e:
            logger.error(f"Error plotting prediction: {e}")
    
    def check_and_retrain(self):
        """
        Check if model retraining is needed and retrain if necessary.
        
        Returns:
        --------
        bool
            True if retrained, False otherwise
        """
        if self.last_training_time is None:
            logger.info("No previous training record. Training model...")
            self.train_model(force=True)
            return True
        
        # Calculate hours since last training
        hours_since_training = (datetime.now() - self.last_training_time).total_seconds() / 3600
        
        if hours_since_training >= self.retraining_hours:
            logger.info(f"{hours_since_training:.1f} hours since last training. Retraining model...")
            self.train_model(force=True)
            return True
        else:
            logger.info(f"Model retraining not needed yet. {hours_since_training:.1f}/{self.retraining_hours} hours since last training.")
            return False


def generate_improved_config():
    """
    Generate an improved configuration for the forex trading bot.
    
    Returns:
    --------
    dict
        Configuration dictionary
    """
    config = {
        "mt5_credentials": {
            "login": 42720146,
            "password": "[REMOVED]",
            "server": "AdmiralsGroup-Demo",
            "path": "C:\\Program Files\\Admirals Group MT5 Terminal\\terminal64.exe"
        },
        
        # Symbol and timeframe settings
        "symbol": "EURUSD",
        "timeframe": "H1",           # Changed from M30 to H1
        "look_back": 120,            # Increased from 60 to 120
        "retraining_hours": 24,      # Kept the same
        
        # Trading parameters
        "lot_size": 0.01,            # Kept the same
        "tp_multiplier": 2.5,        # Increased from 2.0 to 2.5
        "sl_multiplier": 1.0,        # Kept the same
        "trailing_start_pct": 0.4,   # Kept the same
        "trailing_step_pct": 0.1,    # Kept the same
        "risk_per_trade_pct": 1.0,   # Kept the same
        
        # Notification settings
        "telegram_bot_token": "[REMOVED]",
        "telegram_chat_id": "[REMOVED]",
        
        # Model parameters
        "confidence_threshold": 0.75,  # Increased from 0.7 to 0.75
        "price_change_threshold": 0.4, # Increased from 0.3 to 0.4
        
        # Data settings
        "max_data_points": 50000      # Increased from 20000 to 50000
    }
    
    return config