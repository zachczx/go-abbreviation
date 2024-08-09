/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./templates/*.{templ,txt}'],
	theme: {
		extend: {
			animation: {
				flash: "flashText 1s ease-out",
			},
			keyframes: {
				flashText: {
					"0%": {
						color: "#44EEFF",
						opacity: 0.1,
					},
					"100%": {color: "#394e6a"},
				},
			},
		},
	},
	plugins: [require('daisyui')],
	daisyui: {
		themes: [
			'light',
			'dark',
			'cupcake',
			'bumblebee',
			'emerald',
			'corporate',
			'synthwave',
			'retro',
			'cyberpunk',
			'valentine',
			'halloween',
			'garden',
			'forest',
			'aqua',
			'lofi',
			'pastel',
			'fantasy',
			'wireframe',
			'black',
			'luxury',
			'dracula',
			'cmyk',
			'autumn',
			'business',
			'acid',
			'lemonade',
			'night',
			'coffee',
			'winter',
			'dim',
			'nord',
			'sunset',
		],
	},
};
