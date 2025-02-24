import kagglehub
import pandas as pd
import os
import duckdb

path = kagglehub.dataset_download("heptapod/titanic")

dataframe = pd.read_csv(os.path.join(path, 'train_and_test2.csv'))
print(dataframe)
duckdb.read_csv(os.path.join(path, 'train_and_test2.csv')).show()

