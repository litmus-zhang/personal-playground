# Give a list containing all operations that can be performed on the numer line, return a binary string where 1 represents the operation is possible and 0 represents the operation is not possible.
# Exakmple 1
# operations = [[1,2], [1, 6], [2, 4, 2], [2, 5, 3]]
# output = "10"

# Example 2
# operations = [[1, 5], [1, 6], [2, 4, 2], [2, 5, 3]]
# output = "00"

# Example 3
# operations = [[1, 5]]
# output = ""


def solution(operations):
    result = ""
    obstacles = []
    for _, operation in enumerate(operations):
        if len(operation) == 2:
            obstacles.append(operation[1])
            continue
        x, size = operation[1], operation[2] - 1
        curr = list(range(x - size, x + size))
        print(curr)
        for i in curr:
            if i in obstacles:
                result += "0"
                break
        else:
            result += "1"
    return result


print(solution([[1, 2], [1, 6], [2, 4, 2], [2, 5, 3], [2, 1, 2], [2, 3, 2]]))