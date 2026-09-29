export type AtLeastTwo<T> = readonly [T, T, ...T[]];

export type TopTwo = readonly [first: number, second: number];

export function topTwo(values: AtLeastTwo<number>): TopTwo;
export function topTwo(values: readonly number[]): TopTwo | undefined;
export function topTwo(values: readonly number[]): TopTwo | undefined {
  if (values.length < 2) {
    return undefined;
  }
  let first = Number.NaN;
  let second = Number.NaN;
  for (const value of values) {
    if (value > first || Number.isNaN(first)) {
      second = first;
      first = value;
    } else if (value > second || Number.isNaN(second)) {
      second = value;
    }
  }
  return [first, second];
}
