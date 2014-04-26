import QtQuick 2.0;

Canvas {
    width: 640
    height: 480
    id: canvas
    property color strokeStyle:  Qt.darker(fillStyle, 1.4)
    property color fillStyle: "#b40000" // red
    onPaint: relations.paintRelations(canvas)
}
