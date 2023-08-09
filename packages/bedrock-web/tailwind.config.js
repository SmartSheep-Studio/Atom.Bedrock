/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx,vue}"],
  theme: {
    extend: {},
    container: {
      center: true,
    },
  },
  plugins: [],
  corePlugins: {
    preflight: false,
  },
};
