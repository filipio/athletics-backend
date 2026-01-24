package utils

import "time"

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
const SessionIDContextKey = ContextKey(5)

const DefaultPageSize = 20
const DefaultPageNumber = 1
const DefaultOrderBy = "id"

const AccessTokenExpiration = time.Hour            // 1 hour
const RefreshTokenExpiration = 90 * 24 * time.Hour // 90 days

type AnyMap map[string]any
