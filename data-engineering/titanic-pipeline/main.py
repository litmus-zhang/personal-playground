import streamlit as st
from data import dataframe


st.title("Titanic data set")
st.write("Displaying the titanic dataset")

st.dataframe(dataframe)

filteredData = dataframe['Parent_Education_Level'].unique()


st.write("Filtered data")
st.dataframe(filteredData)