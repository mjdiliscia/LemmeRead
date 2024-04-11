package ui

var MainWindowUI string = `
<interface>
  <object class="GtkApplicationWindow" id="window">
    <property name="visible">TRUE</property>
  </object>
</interface>
`

var PostUI string = `
<interface>
      <object class="GtkBox" id="post">
        <property name="visible">True</property>
        <property name="can-focus">False</property>
        <property name="orientation">vertical</property>
        <child>
          <object class="GtkLabel" id="title">
            <property name="visible">True</property>
            <property name="can-focus">False</property>
            <property name="label" translatable="yes">Title of post</property>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkBox" id="subtitle">
            <property name="visible">True</property>
            <property name="can-focus">False</property>
            <child>
              <object class="GtkImage" id="communityIcon">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
                <property name="stock">gtk-dialog-warning</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkSeparator">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="padding">3</property>
                <property name="position">1</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="communityName">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
                <property name="label" translatable="yes">Community name</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">2</property>
              </packing>
            </child>
            <child>
              <object class="GtkSeparator">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="padding">3</property>
                <property name="position">3</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="username">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
                <property name="label" translatable="yes">Username</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">4</property>
              </packing>
            </child>
            <child>
              <object class="GtkSeparator">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="padding">3</property>
                <property name="position">5</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="time">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
                <property name="label" translatable="yes">Time</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">6</property>
              </packing>
            </child>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkImage" id="image">
            <property name="visible">True</property>
            <property name="can-focus">False</property>
            <property name="stock">gtk-select-color</property>
          </object>
          <packing>
            <property name="expand">True</property>
            <property name="fill">True</property>
            <property name="position">2</property>
          </packing>
        </child>
        <child>
          <object class="GtkLabel" id="description">
            <property name="visible">True</property>
            <property name="can-focus">False</property>
            <property name="label" translatable="yes">Description of post</property>
            <property name="xalign">0</property>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">3</property>
          </packing>
        </child>
        <child>
          <object class="GtkBox">
            <property name="visible">True</property>
            <property name="can-focus">False</property>
            <child>
              <object class="GtkSpinButton" id="votes">
                <property name="visible">True</property>
                <property name="can-focus">True</property>
                <property name="numeric">True</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkSeparator">
                <property name="visible">True</property>
                <property name="can-focus">False</property>
              </object>
              <packing>
                <property name="expand">True</property>
                <property name="fill">True</property>
                <property name="position">1</property>
              </packing>
            </child>
            <child>
              <object class="GtkButton" id="comments">
                <property name="label" translatable="yes">comments</property>
                <property name="visible">True</property>
                <property name="can-focus">True</property>
                <property name="receives-default">True</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">2</property>
              </packing>
            </child>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">4</property>
          </packing>
        </child>
      </object>
</interface>
`
