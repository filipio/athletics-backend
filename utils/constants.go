package utils

const AdminRole = "admin"
const UserRole = "user"
const OrganizerRole = "organizer"

const OnlyCurrentUserPath = "/users/me"

type ContextKey uint

const UserContextKey = ContextKey(0)
const DbContextKey = ContextKey(1)
const WorkersContextKey = ContextKey(2)
const RecordIdContextKey = ContextKey(3)
const OnlyCurrentUserContextKey = ContextKey(4)

const DefaultPageSize = 20
const DefaultPageNumber = 1
const DefaultOrderBy = "id"

type AnyMap map[string]any
