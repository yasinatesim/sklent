import { create } from "zustand";

import { MODAL, type ModalKind } from "@/shared/constants/modal";

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
