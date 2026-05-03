import styles from './footer.module.css';

export const Footer = () => {
  return (
    <footer className={styles.wrapper}>
      <div className="container">
        <p className={styles.p}>
          <span>
            AI commits powered by{' '}
            <a href="https://diny.run/" target="_blank">
              diny
            </a>
          </span>
        </p>
      </div>
    </footer>
  );
};
