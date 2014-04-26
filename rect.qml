import QtQuick 2.0

Rectangle {
	//anchors.centerIn: parent
    Text {
        text: service.name
    }
	width: 100
	height: 100
	color: "orange"
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        drag.target: parent
    }
}
