// @ts-check

/**
 * @typedef {Object} Hasher
 * @property {(data: Uint8Array) => Promise<Uint8Array>} hash
 * @property {(hash: Uint8Array, data: Uint8Array) => Promise<boolean>} verify
 */

/**
 * @param {Uint8Array} a
 * @param {Uint8Array} b
 * @returns {boolean}
 */
export function timingSafeEqual(a, b) {
  if (!(a instanceof Uint8Array) || !(b instanceof Uint8Array)) {
    throw new TypeError("arguments must be Uint8Array");
  }
  if (a.length !== b.length) {
    return false;
  }
  let mismatch = 0;
  for (let i = 0; i < a.length; i++) {
    mismatch |= a[i] ^ b[i];
  }
  return mismatch === 0;
}
