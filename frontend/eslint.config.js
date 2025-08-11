import tseslint from 'typescript-eslint';
import reactX from 'eslint-plugin-react-x';
import reactDom from 'eslint-plugin-react-dom';
import pluginQuery from '@tanstack/eslint-plugin-query'

export default tseslint.config(
    tseslint.configs.strictTypeChecked,
    tseslint.configs.stylisticTypeChecked,
    tseslint.configs.recommendedTypeChecked,
    reactX.configs['recommended-typescript'],
    reactDom.configs.recommended,
    pluginQuery.configs['flat/recommended'],
    {
        languageOptions: {
            parserOptions: {
                projectService: true,
                project: ['./tsconfig.json'],
            },
        },
        ignores: ["src/clients"]
    },
);
