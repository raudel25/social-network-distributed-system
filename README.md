# social-network-distributed-system

| **Nombre**              | **Grupo** | **Github**                                     |
|-------------------------|-----------|------------------------------------------------|
| Anabel Benítez González | C411      | [@anabel02](https://github.com/anabel02)       |
| Raudel Gómez Molina     | C411      | [@raudel25](https://github.com/raudel25)   |      



## Descripción
Este proyecto busca crear una plataforma de comunicación descentralizada, inspirada en Twitter pero con un enfoque en la privacidad y la resistencia a fallos. Este sistema permite a los usuarios compartir mensajes cortos, seguir a otros y republicar contenido, todo ello en una arquitectura distribuida que garantiza la escalabilidad geográfica y la tolerancia a desconexiones.

## ¿Cómo ejecutarlo?
1. Clona el repositorio
2. Navega al directorio del proyecto:
   ```bash
   cd social-network-distributed-system
   ```
### Local:
3. Instala las dependencias:
   ```bash
   go mod tidy
   ```
4. Inicia el servidor:
    ```bash
   make dev ID=<id>
   ```

### Docker:
3. Construye la imagen de Docker:
   ```bash
   make docker-build
   ```
4. Ejecuta el contenedor:
   ```bash
   make docker-run ID=<id>
   ```
