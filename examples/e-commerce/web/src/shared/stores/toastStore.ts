import { create } from "zustand";

type ToastState = {
  message: string;
  visible: boolean;
  show: (message: string) => void;
  hide: () => void;
};

let timer: ReturnType<typeof setTimeout> | undefined;

export const useToastStore = create<ToastState>((set) => ({
  message: "",
  visible: false,
  show: (message) => {
    set({ message, visible: true });
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => set({ visible: false }), 2500);
  },
  hide: () => set({ visible: false }),
}));
