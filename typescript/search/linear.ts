export function linearSearch<T>(values: readonly T[], target: T): number {
  return values.indexOf(target);
}
