## Steps

  1. Apply a simple diff algorithm on the first level treating each tag and its attributes and its children as a single atom.
  2. Annotate the atrribute changes to each new atom when when there is an existing old atom in the same place. (skip if there were no differences)
  3. Reapply these steps in each atom that was changed from one file to another.

```xml
<!-- Old file -->
<root id="root1" title="Old Title">
    <p class="blue" height="auto">Lorem Ipsum Dolor Est</p>
</root>
```
```xml
<!-- New file -->
<root id="root2" title="New Title">
    <span id="S1" width="100px" height="auto">§1: Lorem Ipsum DOLOR Est</span>
</root>
```

```xml
<!-- Step 1, level root -->
<delta:diff old="ex4_a.xml" new="ex4_b.xml" xmlns:delta="LibreDeltaML">
    <delta:del>
<root id="root1" title="Old Title">
    <p class="blue" height="auto">Lorem Ipsum Dolor Est</p>
</root>
    </delta:del>
    <delta:add>
<root id="root2" title="New Title">
    <span id="S1" width="100px" height="auto">§1: Lorem Ipsum DOLOR Est</span>
</root>
    </delta:add>
</delta:diff>
```

```xml
<!-- Step 2, level root -->
<delta:diff old="ex4_a.xml" new="ex4_b.xml" xmlns:delta="LibreDeltaML">
<root id="root2" title="New Title" delta:old-id="root1" delta:old-title="Old Title" delta:kept="elem" delta:changed="@id|@title">
    <delta:del>
    <p class="blue" height="auto">Lorem Ipsum Dolor Est</p>
    </delta:del>
    <delta:add>
    <span id="S1" width="100px" height="auto">§1: Lorem Ipsum DOLOR Est</span>
    </delta:add>
</root>
</delta:diff>
```

```xml
<!-- Step 1, level root+1 -->
<delta:diff old="ex4_a.xml" new="ex4_b.xml" xmlns:delta="LibreDeltaML">
<root id="root2" title="New Title" delta:old-id="root1" delta:old-title="Old Title" delta:kept="elem" delta:changed="@id|@title">
    <delta:del>
    <p class="blue" height="auto">Lorem Ipsum Dolor Est</p>
    </delta:del>
    <delta:add>
    <span id="S1" width="100px" height="auto">§1: Lorem Ipsum DOLOR Est</span>
    </delta:add>
</root>
</delta:diff>
```

```xml
<!-- Step 2, level root+1 -->
<delta:diff old="ex4_a.xml" new="ex4_b.xml" xmlns:delta="LibreDeltaML">
<root id="root2" title="New Title" delta:old-id="root1" delta:old-title="Old Title" delta:kept="elem" delta:changed="@id|@title">
    <span id="S1" width="100px" height="auto" delta:old-class="blue" delta:kept="@height" delta:was-elem="p" delta:was-changed="elem|@id|@width|@class"><delta:del>Lorem Ipsum Dolor Est</delta:del><delta:add>§1: Lorem Ipsum DOLOR Est</delta:add></span>
</root>
</delta:diff>
```

```xml
<!-- Step 1, level root+2 -->
<delta:diff old="ex4_a.xml" new="ex4_b.xml" xmlns:delta="LibreDeltaML">
<root id="root2" title="New Title" delta:old-id="root1" delta:old-title="Old Title" delta:kept="elem" delta:changed="@id|@title">
    <span id="S1" width="100px" height="auto" delta:old-class="blue" delta:kept="@height" delta:was-elem="p" delta:was-changed="elem|@id|@width|@class"><delta:add>§1: </delta:add>Lorem Ipsum D<delta:del>olor</delta:del><delta:add>OLOR</delta:add> Est</span>
</root>
</delta:diff>
```