import Link from "next/link";
import styles from "./not-found.module.scss";

const NotFound = () => (
  <main className={styles.wrap}>
    <h1 className={styles.code}>404</h1>
    <Link className={styles.link} href="/">
      Vela Commerce
    </Link>
  </main>
);

export default NotFound;
