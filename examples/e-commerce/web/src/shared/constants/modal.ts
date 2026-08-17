export const MODAL = { NONE: "none", COUPON: "coupon" } as const;

export type ModalKind = (typeof MODAL)[keyof typeof MODAL];
