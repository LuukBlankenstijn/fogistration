import { useEffect, useRef, useState } from "react"
import { UserRole, type User } from "@/clients/generated-client"
import { useGetCurrentUser } from "@/query/auth"

interface Props {
  open: boolean
  user: User | null
  onClose: () => void
  onSave: (u: User) => Promise<void> | void
}
export function EditUserModal({ open, user, onClose, onSave }: Props) {
  const [form, setForm] = useState<User | null>(user)
  const first = useRef<HTMLInputElement>(null)
  const { data: currentUser } = useGetCurrentUser()

  useEffect(() => {
    setForm(user)
    if (open) setTimeout(() => first.current?.focus(), 0)
  }, [user, open])

  if (!open || !form) return null

  const set = <K extends keyof User>(k: K, v: User[K]) => {
    setForm({ ...form, [k]: v })
  }

  const isValid =
    form.username.trim().length > 0 &&
    form.email.trim().length > 0

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-black/60">
      <div
        className="w-[520px] rounded-2xl border p-6 shadow-2xl"
        style={{
          backgroundColor: "hsl(var(--panel))",
          borderColor: "hsl(var(--border))",
        }}
      >
        <h3
          className="mb-4 text-lg font-semibold"
          style={{ color: "hsl(var(--fg))" }}
        >
          Edit user
        </h3>
        <form className="space-y-4">
          <div>
            <label
              className="mb-1 block text-sm"
              style={{ color: "hsl(var(--muted))" }}
            >
              Name
            </label>
            <input
              ref={first}
              value={form.username}
              onChange={(e) => {
                set("username", e.target.value)
              }}
              className="w-full rounded-md border px-3 py-2"
              style={{
                backgroundColor: "hsl(var(--panel))",
                borderColor: "hsl(var(--border))",
                color: "hsl(var(--fg))",
              }}
            />
          </div>
          <div>
            <label
              className="mb-1 block text-sm"
              style={{ color: "hsl(var(--muted))" }}
            >
              Email
            </label>
            <input
              value={form.email}
              onChange={(e) => {
                set("email", e.target.value)
              }}
              className="w-full rounded-md border px-3 py-2"
              style={{
                backgroundColor: "hsl(var(--panel))",
                borderColor: "hsl(var(--border))",
                color: "hsl(var(--fg))",
              }}
            />
          </div>
          <div>
            <label
              className="mb-1 block text-sm"
              style={{ color: "hsl(var(--muted))" }}
            >
              Role
            </label>
            <select
              value={form.role as string}
              onChange={(e) => {
                set("role", e.target.value as User["role"])
              }}
              className="w-full rounded-md border px-3 py-2
                         disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={form.id === currentUser.id}
              style={{
                backgroundColor: "hsl(var(--panel))",
                borderColor: "hsl(var(--border))",
                color: "hsl(var(--fg))",
              }}
            >
              <option value={UserRole.ADMIN}>Admin</option>
              <option value={UserRole.USER}>User</option>
              <option value={UserRole.GUEST}>Guest</option>
            </select>
          </div>

          <div className="mt-6 flex justify-end gap-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-md border px-3 py-2 transition-colors hover:bg-[hsl(var(--hover))]"
              style={{
                borderColor: "hsl(var(--border))",
                color: "hsl(var(--fg))",
              }}
            >
              Cancel
            </button>
            <button
              type="button"
              disabled={!isValid}
              onClick={() => void onSave(form)}
              className="
    rounded-md border px-3 py-2 text-[hsl(var(--fg))]
    transition-colors hover:bg-[hsl(var(--hover))]
    disabled:text-[hsl(var(--muted))]
    disabled:border-[hsl(var(--border))]
    disabled:bg-transparent
    disabled:hover:bg-transparent
    disabled:cursor-not-allowed
  "
            >
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
