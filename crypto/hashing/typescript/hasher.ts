export interface Hasher {
  hash(data: Uint8Array): Promise<Uint8Array>;
  verify(hash: Uint8Array, data: Uint8Array): Promise<boolean>;
}

function timingSafeEqual(a: Uint8Array, b: Uint8Array): boolean {
  if (a.length !== b.length) {
    return false;
  }
  let mismatch = 0;
  for (let i = 0; i < a.length; i++) {
    mismatch |= a[i] ^ b[i];
  }
  return mismatch === 0;
}
