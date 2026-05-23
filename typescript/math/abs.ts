export function abs(x: number): number;
export function abs(x: bigint): bigint;
export function abs(x: number | bigint): number | bigint {
  return typeof x === "bigint" ? (x < 0n ? -x : x) : Math.abs(x);
}
