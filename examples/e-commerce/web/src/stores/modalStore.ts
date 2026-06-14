import { create } from "zustand";

export const MODAL = { NONE: "none", COUPON: "coupon" } as const;
export type ModalKind = (typeof MODAL)[keyof typeof MODAL];

type ModalState = {
  active: ModalKind;
  open: (kind: ModalKind) => void;
  close: () => void;
};

export const useModalStore = create<ModalState>((set) => ({
  active: MODAL.NONE,
  open: (kind) => set({ active: kind }),
  close: () => set({ active: MODAL.NONE }),
}));
