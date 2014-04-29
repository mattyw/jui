import QtQuick 2.0

Rectangle {
    Text {
        anchors.centerIn: parent
        text: service.name
    }
	width: 100
	height: 100
	color: "orange"
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        drag.target: parent
        onReleased: service.newPos(parent.x + 40, 450 - parent.y)
    }
}
