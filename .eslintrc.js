module.exports = {
	"parser": "babel-eslint",
	"plugins": ["react"],
	"env": {
		"browser": true
	},
	"extends": ["eslint:recommended", "plugin:react/recommended"],
	"rules": {
		"max-len": [1, 120, 2, {ignoreComments: true}],
		"react/prop-types": "off",
		"no-console": "off"
	}
};