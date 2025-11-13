# Guía de Usuario - Users Service API

## 📖 Descripción
Microservicio de usuarios con API REST y GraphQL, construido con Go, Gin, GORM y PostgreSQL. Incluye autenticación JWT, documentación Swagger y playground de GraphQL.

## 🚀 Inicio Rápido

### Prerrequisitos
- Docker
- Docker Compose

### Ejecutar el Servicio
```bash
# Clonar el repositorio y navegar al directorio
cd users-api

# Ejecutar con Docker Compose
docker compose up
```

El servicio estará disponible en `http://localhost:3001`

## 📊 Endpoints Disponibles

### 🔐 Autenticación (Públicos)

#### 1. Registrar Usuario
```http
POST /api/auth/register
Content-Type: application/json

{
  "firstName": "Juan",
  "lastName": "Pérez",
  "email": "juan@example.com",
  "username": "juanperez",
  "password": "miContraseñaSegura123"
}
```

**Respuesta:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "username": "juanperez",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. Iniciar Sesión
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "juanperez",
  "password": "miContraseñaSegura123"
}
```

**Respuesta:** Establece una cookie `auth_token` HTTP-only y devuelve:
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "username": "juanperez"
  }
}
```

#### 3. Cerrar Sesión
```http
POST /api/auth/logout
```

### 👥 Gestión de Usuarios (Protegidos - Requieren Autenticación)

#### 4. Obtener Todos los Usuarios
```http
GET /api/users
Cookie: auth_token=<jwt-token>
```

**Respuesta:**
```json
[
  {
    "id": 1,
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "username": "juanperez",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### 5. Obtener Usuario por ID
```http
GET /api/users/1
Cookie: auth_token=<jwt-token>
```

#### 6. Actualizar Usuario
```http
PUT /api/users/1
Content-Type: application/json
Cookie: auth_token=<jwt-token>

{
  "firstName": "Juan Carlos",
  "email": "juancarlos@example.com"
}
```

#### 7. Eliminar Usuario (Soft Delete)
```http
DELETE /api/users/1
Cookie: auth_token=<jwt-token>
```

## 🎮 GraphQL

### Endpoint GraphQL
```http
POST /graphql
Cookie: auth_token=<jwt-token>
```

### Playground GraphQL (Solo desarrollo)
```http
GET /playground
```
Disponible en: `http://localhost:3001/playground`

**Ejemplo de consulta GraphQL:**
```graphql
query {
  users {
    id
    firstName
    lastName
    email
    username
  }
  
  user(id: 1) {
    id
    firstName
    lastName
    email
  }
}
```

## 📚 Documentación Swagger

La documentación interactiva de la API está disponible en:
```http
GET /swagger/index.html
```
URL: `http://localhost:3001/swagger/index.html`

## 🛠️ Health Check

Verificar el estado del servicio:
```http
GET /health
```

**Respuesta:**
```json
{
  "status": "OK",
  "service": "users-service"
}
```

## 🔧 Configuración

### Variables de Entorno
El servicio usa las siguientes variables (configuradas en docker-compose):

- `DATABASE_URL`: Conexión a PostgreSQL
- `JWT_SECRET`: Clave secreta para JWT
- `PORT`: Puerto del servicio (3001)
- `ENVIRONMENT`: Entorno (development/production)

### Base de Datos
- **Automático**: GORM crea automáticamente las tablas al iniciar
- **PostgreSQL**: Usa la imagen `postgres:15-alpine`
- **Puerto**: 5432 (accesible localmente para debugging)

## 🐳 Docker

### Estructura de Dockerfiles
- `Dockerfile`: Para binario pre-compilado (producción)
- `Dockerfile.dev`: Para desarrollo con generación de código

### Comandos Útiles
```bash
# Construir y ejecutar
docker compose up

# Solo construir
docker compose build

# Ejecutar en segundo plano
docker compose up -d

# Ver logs
docker compose logs -f

# Detener servicios
docker compose down
```

## 🔒 Autenticación

### Flujo de Autenticación
1. **Registro/Login**: Obtener cookie `auth_token`
2. **Requests Protegidos**: Incluir cookie automáticamente
3. **Logout**: Eliminar cookie

### Características de Seguridad
- ✅ Cookies HTTP-only
- ✅ JWT con expiración (24 horas)
- ✅ Validación de contraseñas hash
- ✅ CORS configurado
- ✅ SameSite cookies

## 🌐 Accesos Directos

Una vez ejecutado el servicio, accede a:

- **API Principal**: `http://localhost:3001`
- **Swagger UI**: `http://localhost:3001/swagger/index.html`
- **GraphQL Playground**: `http://localhost:3001/playground`
- **Health Check**: `http://localhost:3001/health`
- **PostgreSQL**: `localhost:5432` (usuario: myuser, BD: usersdb)

## 🐛 Troubleshooting

### Problemas Comunes
1. **Puerto en uso**: Cambiar puerto en docker-compose.yml
2. **Error de base de datos**: Verificar que PostgreSQL esté saludable
3. **Cookie no persiste**: Verificar configuración de CORS en el frontend

### Logs de Docker
```bash
# Ver todos los logs
docker compose logs

# Ver logs específicos de un servicio
docker compose logs users-service

# Seguir logs en tiempo real
docker compose logs -f users-service
```

---

**¡Listo!** Tu microservicio de usuarios está funcionando con API REST, GraphQL, documentación automática y base de datos persistente. 🚀