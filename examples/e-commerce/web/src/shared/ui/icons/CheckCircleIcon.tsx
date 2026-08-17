type CheckCircleIconProps = {
  size?: number;
};

const CheckCircleIcon = ({ size = 72 }: CheckCircleIconProps) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" aria-hidden="true">
    <circle cx="12" cy="12" r="10" />
    <path d="m9 12 2 2 4-4" />
  </svg>
);

export default CheckCircleIcon;
