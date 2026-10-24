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
        admin: {
          50: '#f8fafc',
          100: '#f1f5f9',
          600: '#0f172a',
          800: '#1e293b',
          900: '#0f172a',
        },
        brand: {
          600: '#2563eb',
        },
      },
    },
  },
  plugins: [],
}
