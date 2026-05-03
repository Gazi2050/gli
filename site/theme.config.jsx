export default {
  logo: <span style={{ fontWeight: 600 }}>gli</span>,
  project: {
    link: 'https://github.com/Gazi2050/gli',
  },
  docsRepositoryBase: 'https://github.com/Gazi2050/gli/tree/main/site/design',
  darkMode: false,
  nextThemes: {
    defaultTheme: 'dark',
    forcedTheme: 'dark',
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – gli',
    };
  },
  feedback: {
    content: null,
  },
  footer: {
    text: (
      <span>
        MIT {new Date().getFullYear()} ©{' '}
        <a href="https://github.com/Gazi2050/gli" target="_blank">
          gli
        </a>
        . AI commits powered by{' '}
        <a href="https://diny.run/" target="_blank">
          diny
        </a>
        .
      </span>
    ),
  },
};
