import { Section } from "./Section"
import { UsersTable } from "./UsersTable"

const Settings = () => {
  return (
    <main className="relative mx-auto max-w-6xl px-6 py-8">
      <Section fallback="je moeder is dik">
        <UsersTable />
      </Section>
    </main>
  )
}

export default Settings
