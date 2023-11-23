/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: [
    "./views/**/*.html",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["light", "dark"],
  }
}

