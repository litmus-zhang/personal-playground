def run_backtest(strategy_fn, data):
    positions = []
    for candle in data:
        signal = strategy_fn(candle)
        positions.append(signal)
    return compute_pnl(positions, data)