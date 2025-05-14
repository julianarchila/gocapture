# GoCapture - Controles de la Interfaz de Usuario

Este documento detalla las reglas de navegación del cursor y los controles de teclado para la interfaz de usuario de terminal de GoCapture.

## Controles Globales

Estos controles funcionan en toda la aplicación:

| Tecla     | Acción                            |
|-----------|-----------------------------------|
| `q`       | Salir de la aplicación            |
| `Ctrl+C`  | Salir de la aplicación            |
| `Esc`     | Volver a la pantalla anterior     |

## Navegación del Menú Principal

El menú principal es el punto de partida de la aplicación:

| Tecla     | Acción                            |
|-----------|-----------------------------------|
| `↑` / `k` | Mover cursor hacia arriba         |
| `↓` / `j` | Mover cursor hacia abajo          |
| `Enter`   | Seleccionar opción resaltada      |

Opciones disponibles en el menú:
- **Iniciar Captura**: Comenzar a capturar tramas en la interfaz especificada
- **Cargar Captura**: Explorar y cargar capturas guardadas previamente
- **Salir**: Salir de la aplicación

## Controles de la Pantalla de Captura

Cuando se están capturando tramas activamente:

| Tecla     | Acción                                 |
|-----------|----------------------------------------|
| `Esc`     | Detener captura y volver al menú principal |
| `Enter`   | Detener captura y ver tramas capturadas|
| `s`       | Guardar la captura actual              |

## Navegación de la Lista de Tramas

Al ver la lista de tramas capturadas:

| Tecla     | Acción                               |
|-----------|--------------------------------------|
| `↑` / `k` | Mover cursor hacia arriba            |
| `↓` / `j` | Mover cursor hacia abajo             |
| `PgUp`    | Mover cursor una página arriba       |
| `PgDn`    | Mover cursor una página abajo        |
| `Enter`   | Ver información detallada de la trama seleccionada |
| `s`       | Guardar la lista actual de tramas    |
| `Esc`     | Volver al menú principal             |

La vista de lista de tramas muestra:
- ID de Trama
- Marca de Tiempo
- Direcciones MAC de origen y destino
- Longitud de la trama
- Resumen del tipo de trama y contenido

## Controles de la Vista de Detalles de Trama

Al ver información detallada de una trama específica:

| Tecla         | Acción                                   |
|---------------|------------------------------------------|
| `Tab`         | Cambiar entre modos de vista (Resumen, Detalles, Hex Dump) |
| `↑` / `k`     | Desplazar contenido hacia arriba         |
| `↓` / `j`     | Desplazar contenido hacia abajo          |
| `→` / `l` / `n` | Ver siguiente trama en secuencia     |
| `←` / `h` / `p` | Ver trama anterior en secuencia      |
| `Esc`         | Volver a la lista de tramas              |

### Modos de Vista

1. **Vista de Resumen**
   - Muestra una visión general concisa de la trama
   - Muestra tipo de trama, direcciones y campos importantes
   - Proporciona contexto de análisis e información de seguridad/QoS

2. **Vista de Detalles**
   - Muestra todos los campos de la trama y sus valores
   - Muestra información raw del encabezado
   - Muestra resultados completos del análisis

3. **Vista Hex Dump**
   - Muestra datos binarios raw en formato hexadecimal
   - Muestra tanto valores hex como representación ASCII
   - Muestra desplazamientos de bytes para fácil referencia

## Pantalla de Capturas Guardadas

Al explorar capturas guardadas:

| Tecla     | Acción                               |
|-----------|--------------------------------------|
| `↑` / `k` | Mover cursor hacia arriba            |
| `↓` / `j` | Mover cursor hacia abajo             |
| `Enter`   | Cargar captura seleccionada          |
| `Esc`     | Volver al menú principal             |

### Selección de Interfaz

Los nombres de interfaz distinguen entre mayúsculas y minúsculas y deben coincidir exactamente con los listados al ejecutar GoCapture sin argumentos.
