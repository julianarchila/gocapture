# GoCapture - Documentación de Usuario

## Descripción General

GoCapture es una herramienta completa de captura y análisis de tramas de red para redes IEEE 802.3 (Ethernet) e IEEE 802.11 (WLAN). Proporciona información detallada sobre el tráfico de red, mecanismos de seguridad y parámetros de Calidad de Servicio (QoS), siendo valiosa para administradores de red, analistas de seguridad y estudiantes que aprenden sobre protocolos de red.

## Instalación

### Prerrequisitos

- Go 1.18 o posterior
- Paquete de desarrollo libpcap
  - Debian/Ubuntu: `sudo apt-get install libpcap-dev`
  - CentOS/RHEL: `sudo yum install libpcap-devel`
  - macOS: `brew install libpcap`

### Compilación desde el Código Fuente

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/USERNAME/gocapture.git
   cd gocapture
   ```

2. Compilar el proyecto:
   ```bash
   go build -o gocapture ./cmd/gocapture
   ```

3. Instalar en el sistema (opcional):
   ```bash
   go install ./cmd/gocapture
   ```

## Uso Básico

### Ver Interfaces de Red Disponibles

Para ver todas las interfaces de red disponibles en su sistema:

```bash
./gocapture
```

### Iniciar Captura de Paquetes

Para comenzar a capturar paquetes en una interfaz específica:

```bash
sudo ./gocapture -interface eth0
```

Nota: Se requieren privilegios de administrador/root para capturar paquetes en la mayoría de las interfaces.

### Opciones de Línea de Comandos

- `-interface`: Interfaz de red desde la cual capturar (ej., eth0, wlan0)
- `-promiscuous`: Habilitar modo promiscuo (predeterminado: true)
- `-filter`: Expresión de filtro BPF (ej., "port 80" para capturar solo tráfico HTTP)

Ejemplo con filtro:
```bash
sudo ./gocapture -interface eth0 -filter "port 53"
```

## Arquitectura de la Aplicación

GoCapture sigue una arquitectura modular con clara separación de responsabilidades:

### Componentes

1. **Motor de Captura** (`internal/capture`)
   - Interactúa con el hardware de red usando la biblioteca libpcap
   - Maneja la captura de tramas en modo promiscuo
   - Proporciona una API basada en canales para consumir tramas capturadas

2. **Módulo Parser** (`internal/parser`)
   - Decodifica tramas raw en datos estructurados
   - Implementa parsers específicos para tramas Ethernet y WLAN
   - Extrae campos de encabezado, direcciones e información específica del protocolo

3. **Motor Analizador** (`internal/analyzer`)
   - Interpreta datos de tramas parseadas
   - Proporciona análisis de seguridad para métodos de encriptación
   - Analiza parámetros QoS y priorización de tráfico
   - Da contexto y recomendaciones basadas en las tramas observadas

4. **Módulo de Almacenamiento** (`internal/storage`)
   - Serializa tramas capturadas a disco
   - Carga capturas guardadas previamente
   - Gestiona metadatos de captura

5. **Interfaz de Usuario** (`ui/`)
   - UI basada en terminal construida con Bubble Tea
   - Múltiples vistas: menú principal, captura, lista de tramas, detalles de trama
   - Navegación interactiva e inspección de tramas

### Flujo de Datos

1. El Motor de Captura captura tramas raw desde la interfaz de red
2. Las tramas raw se envían al Módulo Parser para decodificación
3. Las tramas parseadas se pasan al Motor Analizador para interpretación
4. La UI muestra las tramas analizadas y permite la interacción
5. El Módulo de Almacenamiento puede guardar/cargar capturas en cualquier momento

## Tipos de Tramas

### Tramas Ethernet (IEEE 802.3)

Las tramas Ethernet son la base de las redes cableadas e incluyen:

- Encabezado MAC con direcciones de origen y destino
- Campo EtherType indicando el protocolo de carga útil (ej., IPv4, IPv6, ARP)
- Etiquetado VLAN opcional para segmentación de red
- Datos de carga útil
- Secuencia de Verificación de Trama (FCS) para detección de errores

### Tramas WLAN (IEEE 802.11)

Las tramas WLAN son más complejas y se categorizan en tres tipos:

1. **Tramas de Gestión**
   - Establecen y mantienen comunicaciones
   - Ejemplos: Beacons, Solicitudes/Respuestas de Asociación, Tramas de Autenticación
   - Usadas para descubrimiento de red y gestión de conexiones

2. **Tramas de Control**
   - Asisten en la entrega de tramas de datos
   - Ejemplos: Acuses de Recibo (ACK), Solicitud para Enviar (RTS), Listo para Enviar (CTS)
   - Ayudan a gestionar el acceso al medio inalámbrico compartido

3. **Tramas de Datos**
   - Transportan los datos reales de la red
   - Pueden incluir parámetros QoS para priorización de tráfico
   - Pueden estar protegidas por varios métodos de encriptación

## Análisis de Seguridad

GoCapture identifica y analiza métodos de encriptación usados en redes inalámbricas:

1. **WEP (Wired Equivalent Privacy)**
   - Encriptación heredada con graves fallos de seguridad
   - Usa cifrado RC4 con claves estáticas
   - Vulnerable a ataques estadísticos

2. **WPA (Wi-Fi Protected Access)**
   - Usa TKIP (Temporal Key Integrity Protocol)
   - Más fuerte que WEP, pero aún tiene vulnerabilidades
   - Diseñado como solución transitoria

3. **WPA2**
   - Usa CCMP basado en el algoritmo AES
   - Actualmente el estándar de seguridad Wi-Fi más ampliamente desplegado
   - Vulnerable a ataques KRACK (si no está parcheado)

4. **WPA3**
   - Último estándar de seguridad para Wi-Fi
   - Usa SAE (Simultaneous Authentication of Equals)
   - Proporciona secreto hacia adelante y protección contra ataques de diccionario offline

## Análisis QoS

Para tramas con información de Calidad de Servicio, GoCapture analiza:

- **Categorías de Tráfico**: Fondo, Mejor Esfuerzo, Video, Voz
- **Niveles de Prioridad**: 0 (más bajo) a 7 (más alto)
- **Asignación TXOP**: Duraciones de oportunidad de transmisión
- **Políticas ACK**: Cómo se confirman las tramas

## Solución de Problemas

### Problemas de Permisos

Si encuentra errores "Operation not permitted":

1. Asegúrese de ejecutar con privilegios de administrador/root:
   ```bash
   sudo ./gocapture -interface eth0
   ```

2. Verifique que la interfaz existe y está activa:
   ```bash
   ip link show
   ```

### No Se Capturan Tramas

Si la aplicación se ejecuta pero no captura tramas:

1. Verifique si su interfaz está en modo monitor (para capturas WLAN)
2. Intente una interfaz diferente
3. Asegúrese de que hay tráfico real en la interfaz
4. Intente usar la interfaz loopback (`lo`) y genere algo de tráfico local

### Errores de Compilación

Si encuentra errores de compilación relacionados con pcap:

1. Asegúrese de que los paquetes de desarrollo libpcap están instalados
2. Verifique que los módulos Go están correctamente inicializados
3. Ejecute `go mod tidy` para resolver dependencias

## Uso Avanzado

### Capturando Tramas Inalámbricas

Para capturar tramas 802.11 raw, su interfaz inalámbrica debe estar en modo monitor:

1. Ponga su interfaz en modo monitor (puede variar según el SO y el controlador)
2. Inicie GoCapture con la interfaz en modo monitor

### Usando Filtros BPF

Las expresiones Berkeley Packet Filter (BPF) permiten un filtrado preciso de captura:

- `port 80 or port 443`: Capturar tráfico HTTP y HTTPS
- `host 192.168.1.1`: Capturar tráfico hacia/desde un host específico
- `icmp`: Capturar solo paquetes ICMP
- `not port 22`: Excluir tráfico SSH

## Contribuir

¡Las contribuciones a GoCapture son bienvenidas! Por favor, consulte nuestras guías de contribución para más información.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - vea el archivo LICENSE para más detalles. 