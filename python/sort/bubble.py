from collections.abc import MutableSequence


def bubble_sort(values: MutableSequence[int]) -> None:
    n = len(values)
    while n > 1:
        last_swap = 0
        for i in range(n - 1):
            if values[i + 1] < values[i]:
                values[i], values[i + 1] = values[i + 1], values[i]
                last_swap = i + 1
        n = last_swap
