import styles from "./CardPreview.module.scss";

type CardPreviewProps = {
  number: string;
  name: string;
  exp: string;
};

const cardType = (number: string): string => {
  const first = number.replace(/\s/g, "")[0];
  if (first === "4") return "Visa";
  if (first === "5") return "Mastercard";
  if (first === "6") return "Troy";
  return "Kredi Kartı";
};

const CardPreview = ({ number, name, exp }: CardPreviewProps) => (
  <div className={styles.card}>
    <div className={styles.type}>{cardType(number)}</div>
    <div className={styles.number}>{number || "•••• •••• •••• ••••"}</div>
    <div className={styles.name}>{name.toUpperCase() || "KART SAHİBİ"}</div>
    <div className={styles.exp}>{exp || "AA/YY"}</div>
  </div>
);

export default CardPreview;
