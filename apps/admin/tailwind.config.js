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
        // MOA premium identity: deep pine console + champagne gold
        admin: {
          50: '#EBF2EC',
          100: '#D9E7DD',
          600: '#1E5A38',
          800: '#123522',
          900: '#0C2318',
          950: '#08170F',
        },
        brand: {
          500: '#2A7248',
          600: '#1E5A38',
          700: '#17452C',
        },
        accent: {
          100: '#F2EAD3',
          500: '#C0A04E',
        },
      },
      fontFamily: {
        sans: ['Pretendard', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Malgun Gothic', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
