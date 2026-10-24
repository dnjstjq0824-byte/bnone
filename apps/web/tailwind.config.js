/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb', // Primary
          700: '#1d4ed8', // Dark
          900: '#1e3a8a',
        },
        accent: {
          500: '#0d9488',
          600: '#0f766e',
        },
      },
    },
  },
  plugins: [],
}
