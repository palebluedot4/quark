from collections.abc import MutableSequence


def selection_sort(values: MutableSequence[int]) -> None:
    n = len(values)
    for i in range(n - 1):
        min_index = min(range(i, n), key=values.__getitem__)
        if min_index != i:
            values[i], values[min_index] = values[min_index], values[i]
