export function bubbleSort(values: number[]): void {
  let n = values.length;
  while (n > 1) {
    let lastSwap = 0;
    for (let i = 1; i < n; i++) {
      const left = values[i - 1];
      const right = values[i];
      if (left !== undefined && right !== undefined && right < left) {
        values[i - 1] = right;
        values[i] = left;
        lastSwap = i;
      }
    }
    n = lastSwap;
  }
}
