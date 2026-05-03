import styles from './hero.module.css';
import Link from 'next/link';

export const Hero = () => {
  return (
    <div className={styles.wrapper}>
      <h1 className={styles.heading}>gli</h1>
      <p style={{ marginTop: 0, fontSize: 18, textAlign: 'center' }}>
        A lightweight Git wrapper. AI-powered commits, no config needed.
      </p>
      <div className={styles.buttons}>
        <a
          data-primary=""
          className={styles.button}
          href="https://github.com/Gazi2050/gli/releases"
          target="_blank"
        >
          Download
        </a>
        <Link href="/getting-started" className={styles.button}>
          Documentation
        </Link>
      </div>
      <a
        className={styles.link}
        href="https://github.com/Gazi2050/gli"
        target="_blank"
      >
        GitHub
      </a>
    </div>
  );
};
