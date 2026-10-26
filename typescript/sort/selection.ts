export function selectionSort(values: number[]): void {
  for (const [i, current] of values.entries()) {
    let minIndex = i;
    let min = current;
    for (let j = i + 1; j < values.length; j++) {
      const value = values[j];
      if (value !== undefined && value < min) {
        minIndex = j;
        min = value;
      }
    }
    if (minIndex !== i) {
      values[minIndex] = current;
      values[i] = min;
    }
  }
}
