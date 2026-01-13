import ray
from ray.rllib.algorithms.ppo import PPOConfig
import numpy as np
from env import TradingEnv

ray.init()

prices = np.random.normal(100, 5, 1000)

config = (
    PPOConfig()
    .environment(
        env=TradingEnv,
        env_config={"prices": prices}

    )
    .framework("torch")
    .env_runners(num_rollout_workers=2)
    .training(
        gamma=0.99,
        lr=3e-4,
        clip_param=0.2
    )
)

# Train agent
algo = config.build()

for i in range(10):
    result = algo.train()
    print(f"Iteration {i}, reward: {result['episode_reward_mean']}")


# Evaluate Model
env = TradingEnv({"prices": prices})
obs, _ = env.reset()

done = False
while not done:
    action = algo.compute_single_action(obs)
    obs, reward, done, _, _ = env.step(action)

