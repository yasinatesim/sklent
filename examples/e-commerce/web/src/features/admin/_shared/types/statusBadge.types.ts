export const STATUS_TONE = {
  ACTIVE: "active",
  PASSIVE: "passive",
  WARNING: "warning",
  DANGER: "danger",
} as const;

export type StatusTone = (typeof STATUS_TONE)[keyof typeof STATUS_TONE];
