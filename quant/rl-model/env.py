import gymnasium as gym
import numpy as np

class TradingEnv(gym.Env):
    def __init__(self, config):
        self.prices = config["prices"]
        self.initial_balance = 10000
        self.action_space = gym.spaces.Discrete(3)  # 0=Hold, 1=Buy, 2=Sell
        self.observation_space = gym.spaces.Box(
            low=-np.inf, high=np.inf, shape=(3,), dtype=np.float32
        )
        self.reset()

    def reset(self, seed=None, options=None):
        self.balance = self.initial_balance
        self.position = 0
        self.current_step = 0
        return self._get_obs(), {}

    def _get_obs(self):
        price = self.prices[self.current_step]
        return np.array([price, self.balance, self.position], dtype=np.float32)

    def step(self, action):
        price = self.prices[self.current_step]

        if action == 1 and self.balance > price:
            self.position += 1
            self.balance -= price
        elif action == 2 and self.position > 0:
            self.position -= 1
            self.balance += price

        self.current_step += 1
        done = self.current_step >= len(self.prices) - 1

        net_worth = self.balance + self.position * price
        reward = net_worth - self.initial_balance

        return self._get_obs(), reward, done, False, {}
