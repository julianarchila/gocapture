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

## Reglas y Comportamiento del Cursor

El cursor en GoCapture sigue estas reglas consistentes:

1. **Visibilidad**: El cursor siempre es visible como un carácter `>` al inicio del elemento seleccionado

2. **Envoltura**: Los cursores no se envuelven de abajo hacia arriba o de arriba hacia abajo

3. **Paginación**: Cuando una lista excede el área visible:
   - La vista se desplaza automáticamente para mantener el cursor visible
   - Las teclas PgUp/PgDn mueven el cursor una página completa
   - Se indica la posición actual de la página (ej., "Mostrando 1-10 de 50 tramas")

4. **Selección**: El elemento actualmente seleccionado siempre está resaltado con el cursor

5. **Verificación de Límites**: 
   - Presionar arriba cuando está en la parte superior de una lista no tiene efecto
   - Presionar abajo cuando está en la parte inferior de una lista no tiene efecto
   - La navegación siguiente/anterior de tramas se detiene al principio o final de la lista de tramas

## Manejo Especial de Entrada

### Expresiones de Filtro

Al aplicar expresiones de filtro BPF a través de argumentos de línea de comandos:

- Las expresiones deben estar correctamente entre comillas si contienen espacios
- Las expresiones complejas pueden usar operadores booleanos (AND, OR, NOT)
- Ejemplos:
  ```
  -filter "host 192.168.1.1 and port 80"
  -filter "not port 22"
  -filter "ether host 00:11:22:33:44:55"
  ```

### Selección de Interfaz

Los nombres de interfaz distinguen entre mayúsculas y minúsculas y deben coincidir exactamente con los listados al ejecutar GoCapture sin argumentos.

## Consideraciones de Accesibilidad

- Se proporciona navegación estilo Vim (`h`, `j`, `k`, `l`) como alternativa a las teclas de flecha
- El indicador de cursor de alto contraste (`>`) hace que la selección sea claramente visible
- La información de resumen está formateada consistentemente para lectores de pantalla 