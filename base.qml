import QtQuick 2.0

Rectangle {
    width: 640
    height: 480
    Grid {
        anchors.fill: parent
        objectName: "grid"
        spacing: 5
    }

    Canvas {
        id: canvas
        anchors.fill: parent
        property color strokeStyle:  Qt.darker(fillStyle, 1.4)
        property color fillStyle: "#b40000" // red
        onPaint: {
            var ctx = canvas.getContext('2d');
            var originX = 85
            var originY = 75
            ctx.save();
            ctx.fillRect(0, 0, 100, 100);
        }
    }
}
