const POINTS = "30,170 115,140 200,150 285,80 370,50 455,90 540,60";
const DOTS = [
  [30, 170],
  [115, 140],
  [200, 150],
  [285, 80],
  [370, 50],
  [455, 90],
  [540, 60],
];
const DAYS = [
  [10, "Pzt"],
  [95, "Sal"],
  [180, "Çar"],
  [265, "Per"],
  [350, "Cum"],
  [435, "Cmt"],
  [515, "Paz"],
] as const;

const WeeklyOrdersChart = () => (
  <svg viewBox="0 0 600 200" width="100%" role="img" aria-label="Haftalık sipariş grafiği">
    {DAYS.map(([x, label]) => (
      <text key={label} x={x} y={190} fill="var(--muted)" fontSize="11">
        {label}
      </text>
    ))}
    <polyline points={POINTS} fill="none" stroke="var(--accent)" strokeWidth="3" strokeLinecap="round" />
    {DOTS.map(([cx, cy]) => (
      <circle key={`${cx}-${cy}`} cx={cx} cy={cy} r="4" fill="var(--accent)" />
    ))}
  </svg>
);

export default WeeklyOrdersChart;
