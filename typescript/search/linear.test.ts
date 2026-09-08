import { describe, expect, test } from "vitest";

import { linearSearch } from "./linear.ts";

describe("linearSearch", () => {
  test.each([
    { name: "first", values: [3, 1, 4, 1, 5], target: 3, expected: 0 },
    { name: "middle", values: [3, 1, 4, 1, 5], target: 4, expected: 2 },
    { name: "last", values: [3, 1, 4, 1, 5], target: 5, expected: 4 },
    { name: "duplicate", values: [3, 1, 4, 1, 5], target: 1, expected: 1 },
    { name: "absent", values: [3, 1, 4, 1, 5], target: 2, expected: -1 },
    { name: "single", values: [7], target: 7, expected: 0 },
    { name: "empty", values: [], target: 1, expected: -1 },
  ])("returns $expected for $name", ({ values, target, expected }) => {
    const actual = linearSearch(values, target);
    expect(actual).toBe(expected);
  });

  test("finds a string", () => {
    const actual = linearSearch(["a", "c", "b"], "b");
    expect(actual).toBe(2);
  });

  test("does not match NaN", () => {
    const actual = linearSearch([1, Number.NaN, 2], Number.NaN);
    expect(actual).toBe(-1);
  });

  test("matches negative zero with zero", () => {
    const actual = linearSearch([-0], 0);
    expect(actual).toBe(0);
  });

  test("finds an object by reference", () => {
    const target = { x: 3 };
    const actual = linearSearch([{ x: 1 }, target], target);
    expect(actual).toBe(1);
  });

  test("does not find a structurally equal object", () => {
    const actual = linearSearch([{ x: 1 }, { x: 3 }], { x: 3 });
    expect(actual).toBe(-1);
  });
});
