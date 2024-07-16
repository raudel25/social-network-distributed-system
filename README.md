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

1. Asegúrate de tener Go instalado en tu sistema.

3. Instala las dependencias:
   ```
   make deps
   ```

4. Genera el código de los protocol buffers (si es necesario):
   ```
   make proto
   ```

5. Ejecuta la aplicación en modo desarrollo:
   ```
   make dev
   ```
   Esto iniciará la aplicación con los puertos por defecto (10000, 11000, 12000).

6. Para ejecutar con un ID específico (por ejemplo, ID=1):
   ```
   make dev ID=1
   ```
   Esto usará los puertos 10001, 11001, 12001.

### Docker:

1. Asegúrate de tener Docker instalado en tu sistema.

2. Construye la imagen Docker:
   ```
   make docker-build
   ```

3. Ejecuta el contenedor Docker:
   ```
   make docker-run
   ```
   Esto iniciará la aplicación con los puertos por defecto.

4. Para ejecutar con un ID específico (por ejemplo, ID=1):
   ```
   make docker-run ID=1
   ```

5. Si necesitas generar el código de los protocol buffers dentro del contenedor:
   ```
   make docker-proto
   ```

### Notas adicionales
- Puertos utilizados
   - Puerto principal (10000 + ID): Puerto principal de la aplicación para comunicaciones generales.
   - Puerto BL (11000 + ID): Puerto de escucha para broadcast. Aquí la aplicación recibe mensajes broadcast de otros nodos.
   - Puerto BR (12000 + ID): Puerto para realizar solicitudes broadcast. Desde aquí la aplicación envía mensajes broadcast a otros nodos.
- El directorio del proyecto se monta como volumen en el contenedor Docker, permitiendo editar los archivos localmente y ver los cambios reflejados en el contenedor.

## Funcionalidades

- Implementación de **Chord**: la implementación de **Chord** ofrece las garantías de seguridad ante fallas y replicación de datos para un sistema distribuido usando `k` sucesores (`k` es un parámetro configurable).
- Sincronización de tiempo: el sistema usa el algoritmo de **Berkeleys** para la sincronización de tiempo de todos
sus nodos.
- Auto descubrimiento: el sistema cuenta con una implementación de broadcast que permite descubrir de forma automática nodos en la misma red que cumplan el mismo protocolo.
- Elección de líder: para mantener la consistencia el sistema cuenta con un nodo líder el cual se elige mediante un mecanismo de concenso entre todos los integrantes del sistema y es el encargado de mantener el auto descubrimiento,
la sincronización de tiempo y procesar los pedidos externos al sistema.
- Manejo de versiones: el sistema cuenta con un mecanismo de concenso para solucionar las inconsistencias de datos
con diferentes versiones que pueden ocurrir por la desconexión de elementos del sistema.

## Cliente 

Como parte del proyecto se desarrollo un cliente usando **Stremalit**, que cuenta con servicios de cache, navegación offline y auto descubrimiento, puede encontrar dicho proyecto en el siguiente <a href="https://github.com/anabel02/social-network-client.git">enlace</a> o clonar el submódulo de git asociado a este repositorio. 