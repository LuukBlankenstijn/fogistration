import tseslint from 'typescript-eslint'
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'
import pluginQuery from '@tanstack/eslint-plugin-query'
import pluginRouter from '@tanstack/eslint-plugin-router'

export default tseslint.config(
    tseslint.configs.strictTypeChecked,
    tseslint.configs.stylisticTypeChecked,
    tseslint.configs.recommendedTypeChecked,
    reactX.configs['recommended-typescript'],
    reactDom.configs.recommended,
    pluginQuery.configs['flat/recommended'],
    pluginRouter.configs['flat/recommended'],
    {
        languageOptions: {
            parserOptions: {
                projectService: true,
                project: ['./tsconfig.json'],
            },
        },
        rules: {
            '@typescript-eslint/only-throw-error': 'off',
        },
    },
    { ignores: ['./src/clients'] },
)

