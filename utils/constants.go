package utils

const AdminRole = "admin"
const UserRole = "user"
const OrganizerRole = "organizer"

type ContextKey uint

const UserContextKey = ContextKey(0)
const DbContextKey = ContextKey(1)

const DefaultPageSize = 20
const DefaultPageNumber = 1
const DefaultOrderBy = "id"

type AnyMap map[string]any
